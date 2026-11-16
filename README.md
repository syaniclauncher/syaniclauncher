## SyanicLauncher (work in progress)
This is the rewrite of SyanicLauncher in GoLang with Wails from scratch while using Vanilla HTML/CSS/JS for the frontend.

#### todos:
- Build a decent frontend
- Complete Building the systems in the launcher

> library used for core: [syaniclauncher/minecraft-lib-golang](https://github.com/syaniclauncher/minecraft-lib-golang)

#### rules:
- Readable and maintainable code
- Avoid AI slop, can be only used in making frontend designs, that too must also have taste and not be slopped & maintainable

#### why?

<details>
<summary>(preview of old SyanicLauncher frontend)</summary>
<img src=".github/readme-images/old-syaniclauncher.png" alt="Preview of old SyanicLauncher">
</details>

Old SyanicLauncher was written in Python + Electron and was totally unmaintainable with a monolithic codebase.

<details>
<summary>(preview of mess)</summary>

![Old Monolith Code](.github/readme-images/old-monolith.png)


<p align="center"><i>yeah this shi was a mess</i></p>

</details>

There were occasional issues with the launcher not working as expected on other systems, issues with python itself, compilation times, size of the compiled artifact due to bundled electron and python making it unpleasing and unperformative.

GoLang + Wails + Vue is a much better combination for a modern, maintainable launcher without compromising and would also help me get better at GoLang.

### commands:
For development: `wails dev`
For building: `wails build`

Go and required dependencies must be installed.
