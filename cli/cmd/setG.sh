#!/bin/bash
newVar="$(pwd)/run-scripts"
sed -i "/^export PATH=.*$newVar.*/d" ~/.bashrc
echo "export PATH=\$PATH:$newVar" >>~/.bashrc
exec bash
