
import { RethinkAdapter, RethinkAdapterOptions } from "pims-rethinkdb";

import { Channels } from "./models";

export class DatabaseHandler {

	private instance: RethinkAdapter;

	constructor(options: RethinkAdapterOptions) {
		options.models = [
			Channels
		];

		this.instance = new RethinkAdapter(options);
	}
}
