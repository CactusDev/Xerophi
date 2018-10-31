
import { Cache, Bucket, RequestVerb } from "../internal";

export interface EndpointData {
	cache?: typeof Cache,
	bucket?: typeof Bucket,
	authentication?: {[verb: string]: string}
}

export function Endpoint(data: EndpointData): Function {
	return () => {
	};
}
