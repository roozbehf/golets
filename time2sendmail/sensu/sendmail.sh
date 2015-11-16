#!/usr/bin/env bash

EMAIL="$1"
JSON="$(/etc/sensu/plugins/time2sendmail)"
RC=$?

case $RC in
  2)
  echo $JSON | mail -s "Sensu Error Report" $EMAIL
  ;;

  1)
  echo $JSON | mail -s "Sensu Warning Report" $EMAIL
  ;;

  0)
  echo $JSON | mail -s "Sensu Check Resolved" $EMAIL
  ;;

  *)
  echo "All good."
esac
