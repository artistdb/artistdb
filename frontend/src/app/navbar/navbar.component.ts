import { Component, Output, EventEmitter } from '@angular/core';

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.css']
})
export class NavbarComponent {

  @Output() categoryUpdate = new EventEmitter<string>(); 
  categories = [
    "Artists",
    "Events",
    "Locations"
  ];

  constructor() { };
  
  setCategory(c: string) {
    this.categoryUpdate.emit(c);
  };
}
