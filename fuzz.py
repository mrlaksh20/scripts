#!/usr/bin/env python3

import argparse
import requests
import threading
import time
from queue import Queue
from urllib3.exceptions import InsecureRequestWarning

requests.packages.urllib3.disable_warnings(category=InsecureRequestWarning)

# ---------------- CONFIG ---------------- #

queue = Queue()

# ---------------- ARGUMENTS ---------------- #

parser = argparse.ArgumentParser(
    description="Mini Python Fuzzer"
)

parser.add_argument(
    "-u",
    "--url",
    required=True,
    help="Target URL with FUZZ keyword"
)

parser.add_argument(
    "-w",
    "--wordlist",
    required=True,
    help="Wordlist path"
)

parser.add_argument(
    "-t",
    "--threads",
    type=int,
    default=10,
    help="Number of threads"
)

parser.add_argument(
    "-d",
    "--delay",
    type=float,
    default=0,
    help="Delay between requests"
)

parser.add_argument(
    "-fc",
    "--filter-code",
    default="404",
    help="Filter status codes (comma separated)"
)

parser.add_argument(
    "-fs",
    "--filter-size",
    default="",
    help="Filter response sizes"
)

args = parser.parse_args()

# ---------------- FILTERS ---------------- #

filter_codes = set()
filter_sizes = set()

if args.filter_code:
    filter_codes = set(
        int(x.strip()) for x in args.filter_code.split(",")
    )

if args.filter_size:
    filter_sizes = set(
        int(x.strip()) for x in args.filter_size.split(",")
    )

# ---------------- BANNER ---------------- #

print("[*] Warming up connection...")
time.sleep(0.5)
print("[*] Warmup done. Starting fuzzing...\n")

# ---------------- LOAD WORDLIST ---------------- #

try:
    with open(args.wordlist, "r") as f:
        words = [line.strip() for line in f if line.strip()]

except FileNotFoundError:
    print("[-] Wordlist not found")
    exit()

# ---------------- WORKER ---------------- #

lock = threading.Lock()


def worker():

    session = requests.Session()

    while not queue.empty():

        word = queue.get()

        target = args.url.replace("FUZZ", word)

        try:
            response = session.get(
                target,
                timeout=10,
                allow_redirects=False,
                verify=False
            )

            status = response.status_code
            size = len(response.content)
            location = response.headers.get("Location", "")

            # -------- FILTERS -------- #

            if status in filter_codes:
                queue.task_done()
                continue

            if size in filter_sizes:
                queue.task_done()
                continue

            # ------------------------- #

            timestamp = time.strftime("%H:%M:%S")

            with lock:

                if location:
                    print(
                        f"[{timestamp}] "
                        f"{status:<3} - "
                        f"{size:<5}B - "
                        f"/{word:<20} "
                        f"-> {location}"
                    )

                else:
                    print(
                        f"[{timestamp}] "
                        f"{status:<3} - "
                        f"{size:<5}B - "
                        f"/{word}"
                    )

            if args.delay > 0:
                time.sleep(args.delay)

        except:
            pass

        queue.task_done()


# ---------------- START ---------------- #

for word in words:
    queue.put(word)

threads = []

for _ in range(args.threads):

    t = threading.Thread(target=worker)
    t.daemon = True
    t.start()

    threads.append(t)

queue.join()

print("\n[*] Fuzzing finished")
