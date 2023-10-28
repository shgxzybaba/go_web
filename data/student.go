package data

type Student struct {
	Username string
	Id uint16
}

func FetchAllStudents() (students []Student, err error) {
    rows, err := DB.Query("SELECT id, username FROM student;")
    if err != nil {
        return nil, err // Return nil slice and error
    }
    defer rows.Close() // Ensure rows are closed after the function exits

    for rows.Next() {
        var student Student
        if err := rows.Scan(&student.Id, &student.Username); err != nil {
            return nil, err // Return nil slice and error
        }
        students = append(students, student)
    }

    if err := rows.Err(); err != nil {
        return nil, err // Return nil slice and error if there was an error during iteration
    }

    return students, nil // Return the populated slice and no error
}
