#!/bin/bash
scp ./.env root@homeassistant.local:/addons/tg_bot
scp -r ./* root@homeassistant.local:/addons/tg_bot
