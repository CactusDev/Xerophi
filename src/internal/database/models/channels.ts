
import { Model, Column } from "pims";

@Model({
	table: "channels"
})
export class Channels {
	@Column({ primary: true })
	public id: string;

	@Column() public token: string;
	@Column() public platform: string;
	@Column() public authorization: ServiceAuthorization;
}
