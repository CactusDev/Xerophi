
import { Model, Column } from "pims";

@Model({
	table: "quotes"
})
export class Quotes {

	@Column({ primary: true })
	public id: string;

	@Column() public channel: string;
	@Column() public quote: Component[];
	@Column() public quoteId: number;
	@Column() public quoted?: string;
	@Column() public when: string;
	@Column() public enabled: boolean;
}
