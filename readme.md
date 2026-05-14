<div align="center">
    <img src="/images/Rahab_by_Julius_Schnorr_von_Carolsfeld.png" alt="Rahab woodcut" width="320">   
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="images/rachav-dark.svg">
      <source media="(prefers-color-scheme: light)" srcset="images/rachav-light.svg">
      <img alt="Rachav logo" src="images/rachav-light.svg" width="420">
    </picture>
    <br>
    <br>
</div>

# Rachav
easy naive proxy setup and dashboard

It just helps you:
- generate configs
- run naiveproxy
- manage certificates
- use a simple web dashboard
- avoid unnecessary pain

# AI disclaimer
Frontend is heavily vibe-coded.

Some backend/proxy handler code also contains AI-assisted code.

This project is still an MVP and currently focused on: making things work properly first :)

I don’t like AI-generated code, even though I use it sometimes. In the distant plans, I want to remove all AI-written code from this repo.


# So, how to run

Place the `bin` and `naive` files in any directory and run it...

On the first launch, it will automatically generate a `config.yaml` file, change it)


## How it works

When started, `rachav` tries to launch 
`./naive`
binary (NaiveProxy) from the current working directory.
```
Minimal directory structure:
├── naive
└── rachav
```

Where:
- `rachav` → rachav binary
- `naive` → NaiveProxy binary

Make sure both files exist in the same directory and have execution permissions.

If not, run:
```bash
chmod +x ./rachav
chmod +x ./naive
```

You can run it in background using nohup:
```bash
nohup ./rachav &
```
Or you can create a service (just google how to do it, it's not hard).

todo: create auto installation script

```mermaid
%%{init: {
  "theme": "base",
  "themeVariables": {
    "background": "#FFFFFF",
    "primaryColor": "#FFFFFF",
    "primaryTextColor": "#000000",
    "primaryBorderColor": "#000000",
    "lineColor": "#000000"
  }
}}%%
flowchart TD
subgraph "how it work"
direction TD
    request --> panelCheck[has panel prefix?]
    panelCheck -- "no (or panel disabled)" --> auth[proxy auth]
    panelCheck -- "yes" --> panel
  
    auth --"bad"--> reverseProxy["reverse proxy\n setup port in config"]
    auth --"good"--> naiveeProxy["naive"]


    subgraph panel["panel (you can disable it)"]
    direction TD
    panelApi["[prefix]/api/"] --> api1[api]
    panelIndex["[prefix]/"] --> panelApp["panel js app"]
    panelApp -- "user interface" --> uReadThis?["[prefix]/api/"]
    end
end
```
