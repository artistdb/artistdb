import { Injectable } from '@angular/core';
import { of } from 'rxjs';

import { FieldBase } from './field-base';

@Injectable({
  providedIn: 'root'
})
export class FieldService {

  getFields(data: string) { 
    var parsed = JSON.parse(data);
    console.log(parsed.data);
    
    var fields: FieldBase<string>[] = [];

    for(var i in parsed.data) {
      var item = parsed.data[i] 

      var field: FieldBase<string> = {
        key: item.key,
        label: item.key,
        value: '',
        required: item.required,
        controlType: item.controlType,
        type: item.type,
      };

      fields.push(field);
    }

    return of(fields);
  };
}
