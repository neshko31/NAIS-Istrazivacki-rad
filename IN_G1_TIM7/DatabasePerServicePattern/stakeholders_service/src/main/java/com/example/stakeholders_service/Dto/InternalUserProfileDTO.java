package com.example.stakeholders_service.Dto;

import com.example.stakeholders_service.Model.User;

public class InternalUserProfileDTO {
    private String firstName;
    private String lastName;
    private String userName;
    private String email;
    private String profileImage;
    private String biography;
    private String motto;

    public InternalUserProfileDTO(String firstName, String lastName, String userName, String email, String profileImage, String biography, String motto) {
        this.firstName = firstName;
        this.lastName = lastName;
        this.userName = userName;
        this.email = email;
        this.profileImage = profileImage;
        this.biography = biography;
        this.motto = motto;
    }

    public static InternalUserProfileDTO fromUser(User user) {
        return new InternalUserProfileDTO(
                user.getFirstName(),
                user.getLastName(),
                user.getUsername(),
                user.getEmail(),
                user.getProfileImage(),
                user.getBiography(),
                user.getMotto()
        );
    }

    public String getFirstName() {
        return firstName;
    }

    public void setFirstName(String firstName) {
        this.firstName = firstName;
    }

    public String getLastName() {
        return lastName;
    }

    public void setLastName(String lastName) {
        this.lastName = lastName;
    }

    public String getUserName() {
        return userName;
    }

    public void setUserName(String userName) {
        this.userName = userName;
    }

    public String getEmail() {
        return email;
    }

    public void setEmail(String email) {
        this.email = email;
    }

    public String getProfileImage() {
        return profileImage;
    }

    public void setProfileImage(String profileImage) {
        this.profileImage = profileImage;
    }

    public String getBiography() {
        return biography;
    }

    public void setBiography(String biography) {
        this.biography = biography;
    }

    public String getMotto() {
        return motto;
    }

    public void setMotto(String motto) {
        this.motto = motto;
    }
}
