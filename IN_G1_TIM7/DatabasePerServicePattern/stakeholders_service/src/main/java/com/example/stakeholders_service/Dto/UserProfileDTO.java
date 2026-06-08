package com.example.stakeholders_service.Dto;

import com.example.stakeholders_service.Model.Role;
import com.example.stakeholders_service.Model.User;

public class UserProfileDTO {
    private String firstName;
    private String lastName;
    private String profileImage;
    private String biography;
    private String motto;
    private String username;
    private String email;

    //konstruktori
    public UserProfileDTO() {
    }

    public UserProfileDTO(String firstName, String lastName, String profileImage, String biography, String motto, String username, String email) {
        this.firstName = firstName;
        this.lastName = lastName;
        this.profileImage = profileImage;
        this.biography = biography;
        this.motto = motto;
        this.username = username;
        this.email = email;
    }

    //seteri i geteri
    public String getFirstName() { return firstName; }
    public void setFirstName(String firstName) { this.firstName = firstName; }
    public String getLastName() { return lastName; }
    public void setLastName(String lastName) { this.lastName = lastName; }
    public String getProfileImage() { return profileImage; }
    public void setProfileImage(String profileImage) { this.profileImage = profileImage; }
    public String getBiography() { return biography; }
    public void setBiography(String biography) { this.biography = biography; }
    public String getMotto() { return motto; }
    public void setMotto(String motto) { this.motto = motto; }

    public String getUsername() {
        return username;
    }

    public void setUsername(String username) {
        this.username = username;
    }

    public String getEmail() {
        return email;
    }

    public void setEmail(String email) {
        this.email = email;
    }
}