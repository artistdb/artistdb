import { Component, Input} from '@angular/core';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'artistDB';
  category = '';

  changeCategory (c: string) {
    this.category = c;
    console.log(c);
  };

  show(): void {
    console.log(this.category);

  }
}