export namespace main {
	
	export class SendFileResponse {
	    submission_id?: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new SendFileResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.submission_id = source["submission_id"];
	        this.error = source["error"];
	    }
	}
	export class Task {
	    id: string;
	    type: string;
	    label: string;
	
	    static createFrom(source: any = {}) {
	        return new Task(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.type = source["type"];
	        this.label = source["label"];
	    }
	}
	export class StartupData {
	    title: string;
	    languages: model.UILanguage[];
	    tasks: Task[];
	
	    static createFrom(source: any = {}) {
	        return new StartupData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.languages = this.convertValues(source["languages"], model.UILanguage);
	        this.tasks = this.convertValues(source["tasks"], Task);
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

export namespace model {
	
	export class UILanguage {
	    name: string;
	    extension: string;
	
	    static createFrom(source: any = {}) {
	        return new UILanguage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.extension = source["extension"];
	    }
	}

}

