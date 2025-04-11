<?php
if ($_SERVER["REQUEST_METHOD"] == "POST") {
    // Маягтын өгөгдлийг авах ба цэвэрлэх
    $name = trim($_POST['name']);
    $email = trim($_POST['email']);
    $message = trim($_POST['message']);

    // Үндсэн баталгаажуулалт: талбарууд хоосон байхгүй эсэхийг шалгана
    if (empty($name) || empty($email) || empty($message)) {
        echo "Бүх талбаруудыг бөглөнө үү.";
        exit;
    }

    // И-мэйл форматыг шалгах
    if (!filter_var($email, FILTER_VALIDATE_EMAIL)) {
        echo "И-мэйл хаяг буруу байна.";
        exit;
    }

    // Өгөгдлийг обработлох: жишээ нь, и-мэйл илгээх
    $to = "your_email@example.com"; // Өөрийн и-мэйл хаягаа энд оруулна уу
    $subject = "Шинэ холбоо барих маягтын илгээлэл";
    $body = "Нэр: $name\nИ-мэйл: $email\nМессеж: $message";
    $headers = "From: $email";

    if (mail($to, $subject, $body, $headers)) {
        echo "Бидэнтэй холбоо барьсанд баярлалаа. Бид удахгүй хариу илгээх болно.";
    } else {
        echo "Мессеж илгээхэд алдаа гарлаа. Дахин оролдоно уу.";
    }
}
?>
