#!/usr/bin/env bash

# Definitions
Green_font_prefix="\033[32m"
Red_font_prefix="\033[31m"
Font_color_suffix="\033[0m"
Info="${Green_font_prefix}[INFO]${Font_color_suffix}"
Error="${Red_font_prefix}[ERROR]${Font_color_suffix}"
echoerr() { if [[ ${QUIET:-0} -ne 1 ]]; then echo -e "${Error} $@" 1>&2; fi }
echoinfo() { if [[ ${QUIET:-0} -ne 1 ]]; then echo -e "${Info} $@" 1>&2; fi }


#### Main
if [ "${1:0:1}" = '-' ]; then
    set -- suprasched "$@"
fi
if [ "$1" != 'suprasched' ] && [ ! -e "$1" ]; then
    set -- suprasched "$@"
fi

if [ ! -z ${DEBUG} ];then
    if  ! echo "$@" |grep -Eoq "\-v";then
        set -- $@ -v
    fi
elif [ ! -z ${TRACE} ];then
    if  ! echo "$@" |grep -Eoq "\-t";then
        set -- $@ -t
    fi
fi
export PATH=$PATH:`pwd`:$NEWUSERHOME

if [ ! -z ${DEBUG} ] || [ ! -z ${TRACE} ];then
    echoinfo "Image created: $(cat /opt/build_date) $(suprasched --version)"
fi

exec $@
