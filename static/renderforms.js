export function addEmployee(employee){
    /*
        employee fields:
        id, first, last, email, department, departmentID
    */
   const body = document.querySelector('body');

   const settingsElement = document.createElement('div');
   settingsElement.classList.add('settings-container');

   const infoBox = document.createElement('div');
   infoBox.classList.add('info-box');

   const header = document.createElement('header');
   header.textContent = 'New employee';

   infoBox.appendChild(header);

   const form = document.createElement('form');
   form.setAttribute('method','POST')

   const firstName = createFormElement('text','First Name', 'first');
   const lastName = createFormElement('text', 'Last Name', 'last');
   const email = createFormElement('email','Email','email');

   form.appendChild(firstName);
   form.appendChild(lastName);
   form.appendChild(email);

   const dropdownContainer = document.createElement('div');
   dropdownContainer.className = 'dropdown-container';

   const dropdown = document.createElement('div');
   dropdown.className = 'dropdown';

   const button = document.createElement('button');
   button.type = 'button';
   button.id = 'add';
   //We'll change this into the id of the current department
   //For the moment we'll just use a fake department
   //We'll poll the API later; For the time being, these fake fields will be used
   button.textContent = 'Department';
   dropdown.appendChild(button);

   const svg = document.createElementNS('http://www.w3.org/2000/svg','svg')
   svg.setAttribute('xmlns', 'http://www.w3.org/2000/svg');
   svg.setAttribute('viewBox', '0 0 24 24');
   svg.style.width = '23px';
   svg.style.height = '23px';
   svg.id = 'chevron-down'
   
   //creating the path
   const path = document.createElementNS('http://www.w3.org/2000/svg','path')
   path.setAttributeNS(null, 'd', 'M7.41,8.58L12,13.17L16.59,8.58L18,10L12,16L6,10L7.41,8.58Z')
   svg.appendChild(path)

   dropdown.appendChild(svg);

   const dropdownbox = document.createElement('div');
   dropdownbox.className = 'dropdown-box';


   const title = document.createElement('title');
   title.textContent = 'dropdown';
   svg.appendChild(title);



   //This will be changed with a GET request
   const dropDowndepts = [mockDeptFactory(12,'Engineering'),mockDeptFactory(15, 'Finance')];

   dropDowndepts.forEach(dept => {
    const deptDiv = document.createElement('div');
    deptDiv.textContent = dept.DeptName;
    deptDiv.dataset["department-Id"] = dept.Id;

    deptDiv.addEventListener(`click`, ()=>{
        //Swap elements
        //We don't need to swap if the box is default (which is 'Departments')
        if(button.id == 'add'){
            const temp = deptDiv;
            button.dataset['department-Id'] = temp.dataset['department-Id'];
            button.textContent = temp.textContent;
            deptDiv.remove();
            button.id = 'employee';
        }else{
            const temp = {
                Id: button.dataset['department-Id'],
                DeptName: button.textContent,
            };
            //A (very lazy) swapping apporach
            button.textContent = deptDiv.textContent;
            button.dataset['department-Id'] = deptDiv.dataset['department-Id'];
            //Change it to employee for the next swapping approach
            

            deptDiv.textContent = temp.DeptName;
            deptDiv.dataset['department-Id'] = temp.Id;

        }

        
    });
    dropdownbox.appendChild(deptDiv);

   });
   dropdown.appendChild(dropdownbox);

   dropdownContainer.appendChild(dropdown);
   form.appendChild(dropdownContainer);
   form.appendChild(document.createElement('br'));

   const buttons = document.createElement('div');
   buttons.className = 'buttons';

   const submit = document.createElement('button');
   submit.type = 'submit';
   submit.id = 'submit-info';
   submit.textContent = 'Submit';

   //Add POST request here
   submit.addEventListener(`click`, ()=>{})
   buttons.appendChild(submit);

   const cancel = document.createElement('button');
   cancel.type = 'button';
   cancel.id = 'cancel-info';
   cancel.textContent = 'Cancel';
   cancel.addEventListener(`click`,()=>{
    settingsElement.remove();
   })
   buttons.appendChild(cancel);
   form.appendChild(buttons);

   infoBox.appendChild(form);
   settingsElement.appendChild(infoBox);
   settingsElement.classList.add('show');

   body.insertBefore(settingsElement, document.querySelector('.container'))

}

const createFormElement = (type, placeholder, id) => {
    const ele = document.createElement('input');
    ele.setAttribute('type', type);
    ele.setAttribute('placeholder', placeholder);
    ele.id = id;

    return ele;
}

const mockDeptFactory = (id,deptName) => {
    return {
        Id: id,
        DeptName: deptName,
    };
}