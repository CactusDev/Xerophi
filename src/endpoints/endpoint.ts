
import { Web } from "../web";
import { Config } from "../config";

export abstract class AbstractEndpoint {

	constructor(protected web: Web, protected config: Config) {

	}

	public abstract async init(): Promise<void>;
}