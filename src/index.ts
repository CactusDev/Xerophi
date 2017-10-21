
import "reflect-metadata";

import * as nconf from "config";
import { Config } from "./config";
import { ReflectiveInjector } from "@angular/core";

import { Core } from "./core";
import { Web } from "./web";

const injector = ReflectiveInjector.resolveAndCreate([
    {
        provide: Config,
        useValue: nconf
    },
    {
        deps: [Config],
        provide: Web,
        useFactory: (config: Config) => {
            const web = new Web(config);
            return web;
        }
    },
    Core
]);

const app: Core = injector.get(Core);
app.start().catch(console.error);
