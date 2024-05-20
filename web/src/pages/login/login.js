import React, { useState } from "react";

const SimpleForm = () => {
  const [firstName, setFirstName] = useState("");

  const onFirstNameChange = event => {
    setFirstName(event.target.value);
  };

  return (
    <div>
      <input type="text" name="firstName" onChange={onFirstNameChange} />
      <Greetings firstName={firstName} />
    </div>
  );
};

export default SimpleForm;
