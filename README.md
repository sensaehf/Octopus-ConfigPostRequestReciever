# Octopus-ConfigPostRequestReciever
This program is a web server that receives configuration info for Octopus Scanner Configs

[![CI-GO-Build&Test](https://github.com/sensaehf/Octopus-ConfigPostRequestReciever/actions/workflows/CI-GO.yml/badge.svg?branch=main)](https://github.com/sensaehf/Octopus-ConfigPostRequestReciever/actions/workflows/CI-GO.yml)

```mermaid
flowchart TD
    A[Scanned Source] -->|Post from internal| B(external-confkrabbi.vist.is)
    A -->|Post from external| C(internal-confkrabbi.vist.is)
    C --> D[IIS]
    B --> D
    D --> G[Golang Webserver]
    G -->|Validated and parsed data called to cli| AA[Octopus.exe CLI] 
    AA -->|Save config file to windows| AB[Config files] 
```
