#!/bin/bash
sed -i "s|IP|$(hostname --all-ip-addresses | awk '{print $1}')|" onprem-environment
