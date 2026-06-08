package com.example.stakeholders_service.Service;

import com.example.stakeholders_service.Dto.AdminUserProfileDTO;
import com.example.stakeholders_service.Dto.InternalUserProfileDTO;
import com.example.stakeholders_service.Dto.UserProfileDTO;
import com.example.stakeholders_service.Exception.*;
import com.example.stakeholders_service.Model.Role;
import com.example.stakeholders_service.Model.User;
import com.example.stakeholders_service.Repository.UserRepository;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Objects;
import java.util.stream.Collectors;

@Service
public class UserService {

    public final UserRepository userRepository;

    public UserService(UserRepository userRepository) {
        this.userRepository = userRepository;
    }

    public List<AdminUserProfileDTO> getAllProfiles(Long id) {
        User user = userRepository.findByUserId(id).orElseThrow(() -> new UserNotFoundException(id));
        if (user.getRole() != Role.ADMIN) throw new NonAdminProfileAccessException();

        return userRepository.findByRoleNot(Role.ADMIN).stream()
                .map(AdminUserProfileDTO::fromUser)
                .collect(Collectors.toList());
    }

    public AdminUserProfileDTO getProfileForAdmin(Long adminId, Long userId) {
        User adminUser = userRepository.findByUserId(adminId).orElseThrow(() -> new UserNotFoundException(adminId));
        if (adminUser.getRole() != Role.ADMIN) throw new NonAdminProfileAccessException();

        User user = userRepository.findByUserId(userId).orElseThrow(() -> new UserNotFoundException(userId));
        if (user.getRole() == Role.ADMIN || Objects.equals(userId, adminId)) throw new AdminProfileAccessException();

        return AdminUserProfileDTO.fromUser(user);
    }

    public AdminUserProfileDTO blockProfile(Long adminId, Long userId) {
        User adminUser = userRepository.findByUserId(adminId).orElseThrow(() -> new UserNotFoundException(adminId));
        if (adminUser.getRole() != Role.ADMIN) throw new NonAdminProfileBlockingException();

        User userToBlock = userRepository.findByUserId(userId).orElseThrow(() -> new UserNotFoundException(userId));
        if (userToBlock.getRole() == Role.ADMIN) throw new AdminBlocksOtherAdminException();

        userToBlock.setBlocked(!userToBlock.isBlocked());
        userRepository.save(userToBlock);
        return AdminUserProfileDTO.fromUser(userToBlock);
    }

    public UserProfileDTO getProfile(Long id) {
        User user = userRepository.findByUserId(id).orElseThrow(() -> new UserNotFoundException(id));
        if (user.getRole() == Role.ADMIN) throw new AdminProfileAccessException();
        if (user.isBlocked()) throw new BlockedProfileAccessException();

        return new UserProfileDTO(user.getFirstName(), user.getLastName(), user.getProfileImage(),
                user.getBiography(), user.getMotto(), user.getUsername(), user.getEmail());
    }

    public UserProfileDTO modifyProfile(Long id, UserProfileDTO dto) {
        User user = userRepository.findByUserId(id).orElseThrow(() -> new UserNotFoundException(id));
        if (user.getRole() == Role.ADMIN) throw new AdminProfileAccessException();
        if (user.isBlocked()) throw new BlockedProfileAccessException();

        if (dto.getFirstName() != null) user.setFirstName(dto.getFirstName());
        if (dto.getLastName() != null) user.setLastName(dto.getLastName());
        if (dto.getProfileImage() != null) user.setProfileImage(dto.getProfileImage());
        if (dto.getBiography() != null) user.setBiography(dto.getBiography());
        if (dto.getMotto() != null) user.setMotto(dto.getMotto());

        userRepository.save(user);
        return new UserProfileDTO(user.getFirstName(), user.getLastName(), user.getProfileImage(),
                user.getBiography(), user.getMotto(), user.getUsername(), user.getEmail());
    }
}
