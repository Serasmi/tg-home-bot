#!/bin/bash
scp -O .env zoomer@homeassistant.local:/addons/tg_bot
scp -rO ./* zoomer@homeassistant.local:/addons/tg_bot
