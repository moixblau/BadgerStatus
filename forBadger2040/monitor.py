import sys
import select
import badger2040
import time

badger = badger2040.Badger2040()
badger.set_update_speed(badger2040.UPDATE_NORMAL)

# Historical buffers
DISK_HISTORY = []
MAX_HISTORY = 50

# Draw large CPU and RAM bars
def draw_big_bars(cpu_pct, ram_pct, disk_pct):
    badger.set_pen(15)
    badger.clear()  # white background
    badger.set_pen(0)

    # CPU
    badger.text("CPU", 10, 5, 2)
    badger.rectangle(50, 6, int(200*cpu_pct/100), 10)
    badger.text(f"{cpu_pct}%", 260, 4, 1)

    # RAM
    badger.text("RAM", 10, 20, 2)
    badger.rectangle(50, 21, int(200*ram_pct/100), 10)
    badger.text(f"{int(ram_pct)}%", 260, 19, 1)

    # DISK
    badger.text("DISK", 10, 35, 2)
    badger.rectangle(50, 36, int(200*disk_pct/100), 10)
    badger.text(f"{int(disk_pct)}%", 260, 34, 1)


# Draw historical chart
def draw_history():
    start_x = 40
    start_y = 75
    width = 246
    height = 50

    badger.set_pen(15)
    badger.rectangle(start_x, start_y, width, height)  # chart background
    badger.set_pen(0)

    # Reference lines (0%,25%,50%,75%,100%)
    for perc in [0, 25, 50, 75, 100]:
        y = start_y + height - int(height * perc / 100)
        badger.text(f"{perc}%", 0, y-5, 1)  # label on the left
        badger.rectangle(start_x, y, width, 1)  # horizontal line

    # Draw history
    for i, val in enumerate(DISK_HISTORY):
        x = start_x + int(i * width / MAX_HISTORY)
        disk_height = int(val * height / 100)
        # Disk line
        badger.rectangle(x, start_y + height - disk_height, 2, disk_height)

# Update display
def update_display(cpu, ram_used, ram_total, disk_used, disk_total):
    cpu_pct = cpu
    ram_pct = ram_used / ram_total * 100
    disk_pct = disk_used / disk_total * 100

    DISK_HISTORY.append(disk_pct)

    if len(DISK_HISTORY) > MAX_HISTORY:
        DISK_HISTORY.pop(0)

    draw_big_bars(cpu_pct, ram_pct, disk_pct)
    draw_history()
    badger.update()

# Loop receiving data from the PC
poll = select.poll()
poll.register(sys.stdin, select.POLLIN)

while True:
    if poll.poll(100):  # wait 100ms
        line = sys.stdin.readline().strip()
        if line:
            try:
                # We expect format: "cpu=45 ram_used=2000 ram_total=8000"
                parts = dict(p.split("=") for p in line.split())

                cpu = int(parts.get("cpu", 0))

                ram_used = int(parts.get("ram_used", 0))
                ram_total = int(parts.get("ram_total", 1))

                disk_used = int(parts.get("disk_used", 0))
                disk_total = int(parts.get("disk_total", 1))

                update_display(cpu, ram_used, ram_total, disk_used, disk_total)
            except Exception as e:
                # Show a small error message if something goes wrong
                badger.set_pen(0)
                badger.text("Error: "+str(e), 10, 10, 1)
                badger.update()
    time.sleep(0.1)
