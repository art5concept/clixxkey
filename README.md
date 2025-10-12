# ClixxKey Secure CLI Password Manager

#### Video Demo: https://youtu.be/v8uiq69IVAU

#### Description:

##### Introduction:

Hello, Im Felix Lara from Chiriqui Panama welcome to my final cs50 project.

Clixxkey its a command line password manager built in Go, that use a delay based lock using NTP to hide passwords to myself
designed to secure store, hadle, manage and protect credentials in any desktop OS. Without fancy dependencies, compile and ROLL!!!

My goal was build something practical, and bit unique with a feature of makes me proud and delay based lock using NTP "pool.ntp.org"

###### Features

- **Add, List, Delete Credentials**
- **NTP-Based Delay Lock**: because sometimes i spamm some credentials that i want to hide from myself, now i can set time example 5 days and well its gone it prevents brute force spamming access and you cant cheat changing the system clock.
- **Clipboard Integration**: View a password its copied to you clipboard and After 60 secods is wiped out using memory zeroing mechanisms
- **Cross-Platform**: Runs as a simple excutable on linux mac and windows and with terminal tweaks can alternate screen mode to have a polished CLI experience

###### Technical Details

- **Argon2id key derivation**: generates 32 byte key from your master password using 128KB of memory and 4 threads. ensuring difficulty in GPU attacks
- **Adaptative Encryption**: Checks what kinde of flags you have in your computer via GOARCH if is yes it uses AES-256 for hardware encryption if not fall to xchacha20-poly1305 to compatibility propuses
- **Storage**: Credentials are stored in a encryted JSON file at ~/.clixxkey/passwords.json with salt in ~/.clixxkey/salt.json for key derivation

###### Personal Touch

- **NTP Delay Lock**: I wanted to hide to myself passwords if i tried to access them to early the service fetch the current time from an NTP server, and each entry have an UnlockAfter timestamp. If isnt ready it shows a message with the remeaning time. This ensure Secury even if someone gets temporary access

###### Code structure

```
.
├── README.md
├── cmd
│   └── app
│       └── main.go
├── go.mod
├── go.sum
└── internal
    ├── crypto
    │   └── crypto.go
    ├── models
    │   └── password.go
    ├── repository
    │   ├── file
    │   │   └── file.go
    │   └── repository.go
    └── service
        ├── clipboard.go
        ├── password.go
        ├── service.go
        ├── service_common.go
        ├── service_windows.go
        └── termutil.go

9 directories, 14 files

```

Packages organized in crypto, main, file, service, models, repository and it have also interface implementations to support full storage databases but i use a file one based

- It have also many error hadleing with custom erros like ErrPasswordNotReady, ErrPasswordNOT fund to debbung esay and clear
- CLI UX uses bufio.Scanner for input and text/tabwriter to fancy tables and also platform specific processing for clean screen rendering

##### INSTALATION

1. INSTALL GO `Go 1.20+`
2. CLONE THE REPO `git clone https://github.com/art5concept/clixxkey.git`
3. BUILD in the main directory `go build ./cmd/app/main.go`
4. LAUNCH `./clixxkey` enter your master password (hidden input) and youll see in spanish a menu

```
Password Manager
1. Mostrar Contraseñas
2. Agregar Contraseña
3. Borrar Contraseña
4. Salir
Selecciona una opción:

```

1 list credentials with passwords in asteriscs and you can decide press enter to go back or take an id to see if is posible to see the passwords now
2 add new entry you can add a new website, add the username, the password and the next opening in secods, days, months and years since i want to hide to me in long periods of time if you dont like to add the entry you can just press enter and its pass every entry
3 delete by ID it will show the table of passwords and you can select one if you want to delete
0 Exit cleanly

###### Development Process

I set milestones: that was very difficult to follow but everything was done i started learning go because a friend says that if a could make a url shortener in go he will payme some money and then i started learning go almost 6 months ago for that friend challenge so i startend building this clixxkey project with some knowloge en go then well my milestones was
week 1: understanding how could be a easy file structure and what it gonna do then
week 2: creation of a simple CLI menu witha loop and starting of CRUD in plain text then i realize to its better to manage in JSON files so i move on to
week 3: Encryption with AI tools and meny begginers mistakes but since ill use AI was easier to follow the original documentation so i figureout some golang books form internet and i follow a tutorial in freecode camp and everything goes well
week 4 starting thinking in the NTP locking feature
week 5 handleing clipboard and the end

I learning a bit more of Go crypto libraries, NTP clients, and how make cross platform clipboard API and that was a challenge for me so i use IA again to understand the strange text commandos like the clearscreen print.
I owned the architecture especially the NTP lock

###### Acknowledgments

Thanks to David Malan and all the big staff of CS50 to teach me this skills, thanks to my familie specially to my roots and the new AI tools for helping me to tackle all this tasks

###### More About ME

Im a 28 years old guy who studied mechanical engineering from 2016 to 2023 in the Technological University of Panama then i have to quit because something happens in my family now im studying another degree in a different university a smaller one and well i want apply what i learn in the enginerring field all that math a physics flows inside of me and i hope that i can continue improving all what i am. thanks and bye. Adios.
