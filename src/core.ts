
import * as Hapi from "hapi";

import { Injectable } from "@angular/core";
import { Config } from "./config";

import { Web } from "./web";
import { MongoHandler } from "./mongo";

@Injectable()
export class Core {

    constructor(private config: Config, private web: Web) {

    }

    public async start() {
    	const mongo = new MongoHandler(this.config);
    	await mongo.connect();
        this.web.start(mongo);
    }
}
