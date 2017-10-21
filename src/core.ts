
import * as Hapi from "hapi";

import { Injectable } from "@angular/core";
import { Config } from "./config";

import { Web } from "./web";

@Injectable()
export class Core {

    constructor(private config: Config, private web: Web) {

    }

    public async start() {
        this.web.start();
    }
}
