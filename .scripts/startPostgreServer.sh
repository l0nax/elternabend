#!/bin/bash - 
#===============================================================================
#
#          FILE: startPostgreServer.sh
#
#         USAGE: ./startPostgreServer.sh
#
#   DESCRIPTION: This scripts starts an DEVELOP PostgreSQL Server in a Docker
#                container!
#
#       OPTIONS: ---
#  REQUIREMENTS: docker
#          BUGS: ---
#         NOTES: ---
#        AUTHOR: Francesco Emanuel Bennici (l0nax), benniciemanuel78@gmail.com
#  ORGANIZATION: FABMation GmbH
#       CREATED: 01/15/2020 03:46:50 PM
#      REVISION:  001
#===============================================================================

set -o nounset                              # Treat unset variables as an error

docker pull postgres:9.5

docker run --rm -d \
        --name elternabend_db \
        -e POSTGRES_PASSWORD=postgres \
        -p 5432:5432 \
        postgres:9.5
