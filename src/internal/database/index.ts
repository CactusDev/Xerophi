
import { RethinkAdapter, RethinkAdapterOptions } from "pims-rethinkdb";

import { Aliases, Channels, Commands, Configs,
	     Permits, Points, Quotes, Repeats, Users } from "./models";

export class DatabaseHandler {

	private instance: RethinkAdapter;

	constructor(options: RethinkAdapterOptions) {
		options.models = [
			Aliases, Channels, Commands, Configs, Permits, Points,
			Quotes, Repeats, Users
		];

		this.instance = new RethinkAdapter(options);
	}

	public async setup() {
		await this.instance.ensure();
	}
}
