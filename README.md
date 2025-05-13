# ShutterSync

**ShutterSync** is a fast, Go-powered tool for macOS that transfers photos from your camera or SD card to your SSD. It automatically **renames** files based on EXIF capture time and **organizes** them into a tidy, date-based folder hierarchy.

🖱️ Uses a native macOS **folder picker** — no terminal paths needed.

---

## ✨ Features

* 📸 Transfer photos from any source folder (e.g., camera or SD card)
* 🧠 EXIF-based **renaming** — e.g., `21-04-2025_09-46.HEIC`
* 📂 Auto-organizes photos into nested folders:

  ```
  /2025/04/21/
  ```
* ⚙️ Fast and parallelized using worker goroutines
* 🖥️ Native **macOS folder picker** via [`sqweek/dialog`](https://github.com/sqweek/dialog)

---

## 🚀 Getting Started

### Prerequisites

* macOS
* [Go 1.20+](https://go.dev/dl/)

### Installation

```bash
git clone https://github.com/Elliot727/shuttersync.git
cd shuttersync
make build
```

### Run

```bash
make run
```

You'll be prompted to:

1. Select the **source** folder (e.g. SD card or camera)
2. Select the **destination** folder (e.g. SSD or photos directory)

---

## 📂 File Naming & Folder Structure

Files are named using EXIF `DateTimeOriginal`:

```
21-04-2025_09-46.HEIC
```

And stored in folders like:

```

2025/
    └── 04/
        └── 21/
            └── 21-04-2025_09-46.HEIC
```

No more messy photo dumps — just clean, time-sorted folders.

---

## ⚡ Performance

* Multi-threaded with Go worker pools
* Handles loads of files efficiently
* Skips non-image files and missing EXIF metadata gracefully

---

## 🧼 Cleanup

To remove build artifacts:

```bash
make clean
```

---

## 🧪 Development

Want to contribute?

* Fork the repo
* Open an issue or pull request
* Let's make photo transfers smarter!

