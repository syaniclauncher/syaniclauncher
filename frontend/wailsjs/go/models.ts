export namespace app {
	
	export class MicrosoftDeviceLogin {
	    verification_uri: string;
	    user_code: string;
	
	    static createFrom(source: any = {}) {
	        return new MicrosoftDeviceLogin(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.verification_uri = source["verification_uri"];
	        this.user_code = source["user_code"];
	    }
	}

}

export namespace launcher_accounts {
	
	export class PublicAccount {
	    type: string;
	    username: string;
	    uuid: string;
	    expired?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new PublicAccount(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.username = source["username"];
	        this.uuid = source["uuid"];
	        this.expired = source["expired"];
	    }
	}
	export class AccountList {
	    selected?: string;
	    accounts: PublicAccount[];
	
	    static createFrom(source: any = {}) {
	        return new AccountList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.selected = source["selected"];
	        this.accounts = this.convertValues(source["accounts"], PublicAccount);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

