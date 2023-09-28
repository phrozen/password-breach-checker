# Password Breach Checker

The purpose of this repository is a proof of concept to make a clone of the API found at: [haveibeenpwned](https://haveibeenpwned.com/Passwords)

The general guideline is to explore optimization techniques in Go and use the provided password hash file as much as possible to avoid using any database while at the same time having a high performant service.

To download the password hashes file with counts (ordered by hash) go to [haveibeenpwned](https://haveibeenpwned.com/Passwords) and follow the instructions. Alternatively, go directly into the [PwnedPasswordsDownloader repository](https://github.com/HaveIBeenPwned/PwnedPasswordsDownloader), install the tool and run it locally to get the latest version (around 35.55GB with 866'044,561 hash/count pairs at the time of writing).

Create a `data` folder at the root and save your downloaded `.txt` file there, the database package tests expect the name `data/pwned-passwords-sha1-ordered-by-hash-v8.bin`, but name it as you want if not runnign tests. To get the `.bin` file, compile and run `cmd\process` by passing the `.txt` file as an input with the `-f` flag. This will process the text file into binary with constant size per hash(20) + count(4) pair of 24 bytes. With the binary file, the database package can use **Memory Mapping** to perform **Binary Search** on it for a highly efficient, read-only, search engine.

The `cmd/check` utility can be used to quickly check a password against the database.

The `cmd/server` provides a quick API to check passwords against the database by receiving a SHA1 hash and returning a counter of breaches, where `0` means `Not Found` in the database. It includes a minimal frontend to test the API embedded in the binary. It has some flags to pass the binary database file, specify a port, or disable logging when running load tests.

A simple [K6](https://k6.io/) script is provided in the `load.js` file for Load Testing the API.

