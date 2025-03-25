run-firebase-emulators:
	firebase emulators:start --project testing --import=firebase/emulator_data --export-on-exit=firebase/emulator_data

run-local-infra:
	mkdir -p log
	supervisord -n -c supervisord.infra.conf

stop-local-infra:
	supervisorctl -c supervisord.infra.conf shutdown || true
	$(MAKE) -C backend stop-local-infra
