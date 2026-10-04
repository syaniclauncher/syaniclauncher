## SyanicLauncher (work in progress)
This is the rewrite of SyanicLauncher in GoLang with Wails from scratch while using Vanilla HTML/CSS/JS for the frontend.

#### todos:
- Build a decent frontend
- Go Package for Minecraft Go from scratch (pain)
- Support for various Minecraft modloaders

#### rules:
- Readable and maintainable code
- Avoid AI slop in backend, can be only used in making frontend designs

#### preview of old SyanicLauncher:
![Syanic Launcher](.github/readme-images/old-syaniclauncher.png)

#### why?
Old SyanicLauncher was written in Python + Electron and was totally unmaintainable with a monolithic codebase.

There were occasional issues with the launcher not working as expected, issues with python itself, compilation times, size of the compiled artifact due to bundled electron and python making it unpleasing and unperformative.
