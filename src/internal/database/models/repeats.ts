
import { Model, Column } from "pims";

@Model({
	table: "repeats"
})
export class Repeats {

	@Column({ primary: true })
	public id: string;

	@Column() public channel: string;
	@Column() public command: string;
	@Column() public arguments: string;
	@Column() public enabled: boolean;
	@Column() public interval: number;
}
