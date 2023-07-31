#!/bin/bash

wcorr=68  # manual fix for vertical panels
hcorr=26  # manual fix for horizontal panels


tmps=$(LANG=C xrandr|grep -om1 'current.*,')
tmps=${tmps/,}
tmps=${tmps/current }
echo "screen resolution = $tmps pixels"
wscr=${tmps/ x*}
hscr=${tmps/*x }
wter=$(( (wscr-wcorr)/2 ))
hter=$(( (hscr-hcorr)/2 ))
echo "terminal width  = $wter pixels"
echo "terminal height = $hter pixels"

terminator --geometry="${wter}x${hter}+0-0" -x bash run1.sh &
terminator --geometry="${wter}x${hter}+0-0" -x bash run2.sh &
terminator --geometry="${wter}x${hter}+0-0" -x bash run3.sh &
terminator --geometry="${wter}x${hter}+0-0" -x bash run4.sh &