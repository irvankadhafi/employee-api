SHELL:=/bin/bash

ifdef test_run
	TEST_ARGS := -run $(test_run)
endif

migrate_up=go run . migrate --direction=up --step=0
migrate_down=go run . migrate --direction=down --step=0

run:
	go run . server

migrate:
	@if [ "$(DIRECTION)" = "" ] || [ "$(STEP)" = "" ]; then\
    	$(migrate_up);\
	else\
		go run . migrate --direction=$(DIRECTION) --step=$(STEP);\
    fi

docker:
	@ docker-compose up -d --build