export namespace main {
	
	export class PlanSettings {
	    displayAC: number;
	    displayDC: number;
	    diskAC: number;
	    diskDC: number;
	    sleepAC: number;
	    sleepDC: number;
	
	    static createFrom(source: any = {}) {
	        return new PlanSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.displayAC = source["displayAC"];
	        this.displayDC = source["displayDC"];
	        this.diskAC = source["diskAC"];
	        this.diskDC = source["diskDC"];
	        this.sleepAC = source["sleepAC"];
	        this.sleepDC = source["sleepDC"];
	    }
	}
	export class PowerPlan {
	    guid: string;
	    name: string;
	    active: boolean;
	    hidden: boolean;
	
	    static createFrom(source: any = {}) {
	        return new PowerPlan(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.guid = source["guid"];
	        this.name = source["name"];
	        this.active = source["active"];
	        this.hidden = source["hidden"];
	    }
	}

}

