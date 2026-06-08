package com.example.stakeholders_service.Dto;

import com.example.stakeholders_service.Model.User;

public class AdminUserProfileDTO {
    private Long userId;
    private String firstName;
    private String lastName;
    private String profileImage;
    private String biography;
    private String motto;
    private Boolean isBlocked;

    //konstruktor
    public AdminUserProfileDTO(Long userId, String firstName, String lastName, String profileImage, String biography, String motto, Boolean isBlocked){
        this.userId = userId;
        this.firstName = firstName;
        this.lastName = lastName;
        this.profileImage = profileImage;
        this.biography = biography;
        this.motto = motto;
        this.isBlocked = isBlocked;
    }

    //korisna metoda
    public static AdminUserProfileDTO fromUser(User user) {
        return new AdminUserProfileDTO(
                user.getUserId(),
                user.getFirstName(),
                user.getLastName(),
                user.getProfileImage(),
                user.getBiography(),
                user.getMotto(),
                user.isBlocked()
        );
    }

    //seteri i geteri

    public Long getUserId() {
        return userId;
    }

    public void setUserId(Long userId) {
        this.userId = userId;
    }

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
    public Boolean getBlocked() {
        return isBlocked;
    }
    public void setBlocked(Boolean blocked) {
        isBlocked = blocked;
    }
}
