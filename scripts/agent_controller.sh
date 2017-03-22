#!/bin/bash
#kconfig: 35 90 12
# description: Agent Installer Test
########################################################  . /etc/init.d/functions

# Get function from functions library

# Start the service AgentInstaller

start() {
echo $"***********  agent_controller service started ***********"

command="/opt/infraguard/sbin/infraGuardMain"
daemon "nohup $command >/dev/null 2>&1 &"
#$command


}


stop(){
echo "Going to kill process agent_controller"
pkill  agent_controller.sh

}

### main logic ###
case "$1" in
  start)
        start
        ;;
  stop)
        stop
        ;;

status)
        status agent_controller.sh
        ;;

 *)
        echo $"Usage: $0 {start|stop|status}"
        exit 1
esac

exit 0


