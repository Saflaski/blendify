# Blendify

**Blendify** is a work-in-progress web app that combines people’s music tastes from Last.FM into a single, shared catalogue. Like Spotify blend, blendify uses its own taste match formula as well as combines your top shared artists and songs and presents them to you and your friend(s).

---

## Features

- **Taste Blending** – Merge multiple users’ song preferences into one catalogue
- **Cross platform** - Blendify is platform agnostic, and music tastes can be extracted through LastFM or directly.
- **Dynamic Updates** – Automatically adapts as tastes change  
- **Web-Based** – No installation required
- **Minmal user authentication** - Authenticate once via last.fm and you're good go

---
## Screenshots
### Home page
<img width="1302" height="790" alt="blendify_home_page2" src="https://github.com/user-attachments/assets/ec111fcf-c992-437f-a2d6-76bd1c514c82" />

### Blend page

<img width="1297" height="791" alt="blendify_blend_page" src="https://github.com/user-attachments/assets/34096ada-1c45-4f55-9552-0088c9c42655" />

---
## How It Works

1. Users login through a 3rd party music platform like Last.FM, then send a blend invite to someone.
2. Upon invite acceptance, Blendify analyzes overlapping and unique tastes and presents them.
3. Blendify also calculates a match percentage based on artist, songs, and albums that people have listened to.

---

## Tech Stack

- **Frontend:** React, Javascript, Typescript and TailwindCSS
- **Backend:** Golang
- **Storage:** Redis with persistence
- **APIs:** Last.FM
- **Hosting:** Hetzner VPS
- **Containerization:** Docker

---

## Future
- Generate playlists that can be automatically saved to Spotify/LastFM/Youtube etc
- More data about each artist/song and more direct comparisons between users
- Group blends
- Genre specific blends
