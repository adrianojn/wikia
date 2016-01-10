# YGO Wikia

Tool to translate Ygopro database

Parameters:

	-h      view usage
	-lang   output language (default: English)
	-config configuration file (default: config.json)

Modes:

	-update     update the wikia database
	-translate  translate Ygopro database

Known languages:

	en English
	de German (Deutsch)
	es Spanish (Español)
	fr French (Français)
	it Italian (Italiano)
	ja Japanese (日本語)
	ko Korean (한국어)
	pt Portuguese (Português)
	zh Chinese (中文)

# Usage

This example shows how to translate to Spanish.

1. Create a folder to hold the configuration and database files.
   The commands must to be run from inside this folder.
2. Put `config.json` and `cards.cdb` (the English database) inside the folder.
3. Copy `cards.cdb` to `cards-es.cdb`.
   The "-es" sulfix is the language code.
4. Run the command `wikia -update -lang es`.
   This will download the Wikia database.
5. Run the command `wikia -translate -lang es`.

