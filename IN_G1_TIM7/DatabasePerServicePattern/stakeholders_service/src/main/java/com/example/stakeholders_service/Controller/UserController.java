package com.example.stakeholders_service.Controller;

import com.example.stakeholders_service.Dto.AdminUserProfileDTO;
import com.example.stakeholders_service.Dto.InternalUserProfileDTO;
import com.example.stakeholders_service.Dto.UserProfileDTO;
import com.example.stakeholders_service.Service.UserService;
import com.example.stakeholders_service.Utils.JwtUtil;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
public class UserController {

    private final UserService userService;
    private final JwtUtil jwtUtil;

    public UserController(UserService userService, JwtUtil jwtUtil) {
        this.userService = userService;
        this.jwtUtil = jwtUtil;
    }

    @GetMapping("/api/users/all")
    public ResponseEntity<List<AdminUserProfileDTO>> getAllProfiles(@RequestHeader("Authorization") String authHeader) {
        Long userId = extractUserId(authHeader);
        return ResponseEntity.ok(userService.getAllProfiles(userId));
    }

    @GetMapping("/api/users/all/{userId}")
    public ResponseEntity<AdminUserProfileDTO> getProfileForAdmin(@RequestHeader("Authorization") String authHeader, @PathVariable Long userId) {
        Long adminId = extractUserId(authHeader);
        return ResponseEntity.ok(userService.getProfileForAdmin(adminId, userId));
    }

    @PatchMapping("/api/users/block/{userId}")
    public ResponseEntity<AdminUserProfileDTO> blockProfile(@RequestHeader("Authorization") String authHeader, @PathVariable Long userId) {
        Long adminId = extractUserId(authHeader);
        return ResponseEntity.ok(userService.blockProfile(adminId, userId));
    }

    @GetMapping("/api/users/profile")
    public ResponseEntity<UserProfileDTO> getProfile(@RequestHeader("Authorization") String authHeader) {
        Long userId = extractUserId(authHeader);
        return ResponseEntity.ok(userService.getProfile(userId));
    }

    @PatchMapping("/api/users/profile")
    public ResponseEntity<UserProfileDTO> modifyProfile(@RequestHeader("Authorization") String authHeader, @RequestBody UserProfileDTO userProfileDTO) {
        Long userId = extractUserId(authHeader);
        return ResponseEntity.ok(userService.modifyProfile(userId, userProfileDTO));
    }

    private Long extractUserId(String authHeader) {
        return jwtUtil.extractUserId(authHeader.replace("Bearer ", ""));
    }
}
