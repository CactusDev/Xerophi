
import { reflectAnnotations, createAnnotationFactory } from "reflect-annotations";
import { Cache, Bucket, RequestVerb } from "../internal";

import * as Hapi from "hapi";

type EndpointHandler = (request: Hapi.Request) => object;  // TODO: work out a functional reply type

const HANDLER_DESCRIPTION_KEY = "internal:route:description";
const HANDLER_DATA_KEY = "internal:route:data";
const HANDLER_VERB_KEY = "internal:route:verb";

export interface EndpointData {
	cache?: typeof Cache,
	bucket?: typeof Bucket,
}

export interface HandlerData {
	authorization?: string
}

export function Endpoint(data: EndpointData): Function {
	return (target: Function) => {

		const annotatedMethods = reflectAnnotations(target);
		for (let method of annotatedMethods) {

		}
	};
}

export function Describe(description: string): Function {
	return (target: Function) => {
		Reflect.defineMetadata(HANDLER_DESCRIPTION_KEY, description, target);
	}
}

export function Get(handlerData?: HandlerData): Function {
	return (target: Function) => {
		Reflect.defineMetadata(HANDLER_VERB_KEY, RequestVerb.Get, target);
		Reflect.defineMetadata(HANDLER_DATA_KEY, handlerData, target);
	}
}

export function Post(handlerData?: HandlerData): Function {
	return (target: Function) => {
		Reflect.defineMetadata(HANDLER_VERB_KEY, RequestVerb.Post, target);
		Reflect.defineMetadata(HANDLER_DATA_KEY, handlerData, target);
	}
}

export function Patch(handlerData?: HandlerData): Function {
	return (target: Function) => {
		Reflect.defineMetadata(HANDLER_VERB_KEY, RequestVerb.Patch, target);
		Reflect.defineMetadata(HANDLER_DATA_KEY, handlerData, target);
	}
}

export function Delete(handlerData?: HandlerData): Function {
	return (target: Function) => {
		Reflect.defineMetadata(HANDLER_VERB_KEY, RequestVerb.Delete, target);
		Reflect.defineMetadata(HANDLER_DATA_KEY, handlerData, target);
	}
}
