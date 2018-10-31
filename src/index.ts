
import { Injector } from "dependy";
import { Core } from "./core";

const injector = new Injector(
	{
		injects: Core
	}
);

injector.get(Core).start();
