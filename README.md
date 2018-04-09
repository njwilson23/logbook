# logbook (lb)

## Usage

    lb [-t TOPICNAME] [-d DIRNAME]

Opens an editor to a journal entry labelled with the current date

## Configuration

    LBDIR       default directory for journal files
    LBEDITOR    editor command

## Server

List log entries:

    curl localhost:8080/list

Search logs:

    curl localhost:8080/list?find=INTERESTING_TEXT_SNIPPET

