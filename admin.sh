#!/usr/bin/env bash
kill -9 $(ps aux | grep -E '\./api|\./streamserver|\./scheduler' | grep -v grep | awk '{print $2}')
