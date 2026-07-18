import threading
import time

def transaction_a():
    print("Txn A: acquiring lock A...")
    lock_a.acquire()
    print("Txn A: acquired lock A")
    time.sleep(0.1)
    print("Txn A: acquiring lock B...")
    lock_b.acquire()
    print("Txn A: acquired lock B")
    lock_b.release()
    lock_a.release()
    print("Txn A: done")

def transaction_b():
    print("Txn B: acquiring lock B...")
    lock_b.acquire()
    print("Txn B: acquired lock B")
    time.sleep(0.1)
    print("Txn B: acquiring lock A...")
    lock_a.acquire()
    print("Txn B: acquired lock A")
    lock_a.release()
    lock_b.release()
    print("Txn B: done")


lock_a = threading.Lock()
lock_b = threading.Lock()

ta = threading.Thread(target=transaction_a)
tb = threading.Thread(target=transaction_b)

ta.start()
tb.start()

ta.join(timeout=5)
tb.join(timeout=5)

if ta.is_alive() or tb.is_alive():
    print("DEADLOCK DETECTED: threads are stuck!")
else:
    print("Both transactions completed successfully")
