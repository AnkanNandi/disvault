# DisVault

DisVault is a lightweight file management solution that leverages Discord servers to store, organize, and manage files. **Caution: This is an hobby project and very early in dev**, so it's not recommended to upload large volumes of files to your Discord server as it may exceed Discord's limitations. It does not encrypt any uploaded data.
Latest Binaries can be downloaded [HERE](https://github.com/AnkanNandi/disvault/releases)

> [!WARNING] 
> DisVault is a hobby project, meaning it is not focused on the encryption of files
> THE WINDOWS EXE MAY SHOW A FALSE POSITIVE SO YOU MAY NEED TO WHITELIST IT FOR USAGE, See [here](https://go.dev/doc/faq#virus)
> It is my first ever project written in Go also first project I ever finished to a certain degree
> It may contain bugs or missing features. Use with caution, especially when uploading sensitive or large amounts of data.
> DO NOT STORE SENSITIVE DATA
> Do not Abuse Discord's CDN, your files or server may get deleted or worse your account might get banned

## 📜 **Overview**

DisVault allows you to:

- **Upload:** Automatically Split files into manageable chunks and upload them to your Discord server.
- **Download:** Retrieve your files seamlessly by assembling the chunks back together.
- **List:** View your files in an organized manner, grouped for easier access and management.
- **Delete:** Remove files from the server when they're no longer needed.
- **Groups:** Assign files to groups for easier categorization and searching. In the future, groups may be used for a web hierarchy view where each parent group would become the main folder containing files or other groups under it (NOT A PRIORITY).

### ⚡ Quick Start

1. **Clone the Repository**

   ```bash
   git clone https://github.com/AnkanNandi/disvault.git
   cd disvault
   ```

2. **Set Up Your Environment**

   - Ensure you have Go installed (version 1.23).
   - Set up your Discord bot with the appropriate permissions.

3. **Build and Run**

   ```bash
   go build -o disvault
   ./disvault --help
   ```

4. **Upload Files**

   Use the upload command to add files:

   ```bash
   ./disvault upload --file yourfile.txt
   ```

## 📋 **Features**

- **Discord Integration**: Uses Discord channels for file storage.
- **Searchable**: Easily search and filter files using various flags.
- **Categorization**: Group files to keep everything organized.

## 🚀 **Usage**

```bash
disvault [command]

Available Commands:
  delete      Delete files using their IDs
  download    Download files using their IDs
  group       Group command allows you to create, delete, and manage groups within DisVault.
  help        Help about any command
  list        List the uploaded files
  upload      Upload a file by splitting it into chunks and registering it in the database
  version     Print the version number of DisVault
```

## ⚠️ **Caution**

- **Discord Limitations**: Uploading a large number of files or very large files can exceed Discord’s storage limitations and could get your bot rate-limited or banned.
- **Data Integrity**: This is an hobby project, please don't use disvault as main backup.

## 💡 **Contributing**

Please fork the repository and submit pull requests for new features or bug fixes. For major changes, please open an issue first to discuss what you would like to change.

## 📝 **TODO**

- [ ] Add flags to delete all files, files in a certain group
- [ ] Add flags on downloading files
- [ ] Improve error handling and logging.
- [ ] Implement Tests
- [ ] Enhance file search functionality with more filters.
- [ ] Develop a web-based interface for easier file management.
- [ ] Sync On different devices?

## 🛠️ **Built With**

- [Go](https://golang.org/) - Programming Language
- [SQLite](https://sqlite.org/) - Database
- [Discord API](https://discord.com/developers/docs/intro) - For file storage and retrieval

## 📄 **License**

It is licensed under AGPL-3.0 license, see the [LICENSE](LICENSE) file for details.

## ❤️ **Support**

If you find DisVault useful, please give it a ⭐ on GitHub and share it with your friends!