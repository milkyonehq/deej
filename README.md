<p align="center">
    <a href="https://github.com/milkyonehq/deej">
    <img src=".github/assets/Logo_DeeJ.png" width="80" alt="Logo" /></a>
</p>

<h1 align="center">DeeJ</h1>

<p align="center">A DJ for Discord parties</p>

---

A DJ for Discord parties to easily play music.

## ⏩ Getting Started

### ⚙️ Installation

Install Docker and Docker Compose by using the link below
https://docs.docker.com/engine/install

Clone project
```shell
git clone https://github.com/milkyonehq/deej.git
```

### 🏁 Quickstart

Start DeeJ with Docker Compose
```shell
docker compose -f deployments/docker-compose.yaml up -d
```

Stop DeeJ with Docker Compose
```shell
docker compose -f deployments/docker-compose.yaml down
```

### 🛠️ Environment variables

The DeeJ container can be configured with the environment variables below.

| Variable          | Description                                     | Default Value |
|-------------------|-------------------------------------------------|---------------|
| LOG_LEVEL         | The log level for application output.           | info          |
| DISCORD_BOT_TOKEN | The token used to authenticate the Discord bot. | ""            |

### ⌨️️ Available Discord Commands

DeeJ supports the following slash commands to control music playback and manage the queue.

| Command    | Description                                     | Example Usage                                                                          |
|------------|-------------------------------------------------|----------------------------------------------------------------------------------------|
| `/play`    | Play a single track from a search query or URL. | `/play never gonna give you up` or `/play https://www.youtube.com/watch?v=dQw4w9WgXcQ` |
| `/skip`    | Skip the current track.                         | `/skip`                                                                                |
| `/pause`   | Pause the playback.                             | `/pause`                                                                               |
| `/resume`  | Resume the playback.                            | `/resume`                                                                              |
| `/queue`   | Display the queue.                              | `/queue`                                                                               |
| `/clear`   | Clear the queue.                                | `/clear`                                                                               |
| `/volume`  | Get or set the volume of the player (0-100).    | `/volume` or `/volume 75`                                                              |
| `/shuffle` | Shuffle the queue.                              | `/shuffle`                                                                             |


## 💻 Technologies

<img src="https://skillicons.dev/icons?i=docker,go" alt="technologies" />

## ✏️ License

DeeJ is distributed under the [Apache 2.0 License](LICENSE).

## ✍️ Contributors

Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->

<table>
  <tr>
    <td align="center"><a href="https://github.com/lcarneli"><img src="https://avatars.githubusercontent.com/u/25481821?v=4" width="100px;" alt=""/><br /><sub><b>Lorenzo Carneli</b></sub></a><br /><a href="https://github.com/milkyonehq/deej/commits?author=lcarneli" title="Code">💻</a> <a href="#" title="Ideas">🤔</a></td>
  </tr>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->
<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification. Contributions of any kind welcome!

---

> 🚀 Don't forget to put a ⭐️ on our repositories!