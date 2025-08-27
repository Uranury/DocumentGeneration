# Document Generation Microservice

[![Build Status](https://img.shields.io/gitlab/pipeline-status/document-generator?branch=main)](https://gitlab.com/document-generator/-/pipelines)
[![Python 3.11](https://img.shields.io/badge/python-3.11-blue.svg)](https://www.python.org/downloads/)
[![License](https://img.shields.io/badge/license-Proprietary-blue.svg)](https://)

A high-performance microservice for generating banking/financial documents in DOCX, XLSX, and PDF formats. Built with FastAPI and designed for enterprise-scale document processing.

## Features

- **Multi-format Support**: Generate documents in DOCX, XLSX, HTML, and PDF
- **Template Engine**: Jinja2 templating with dynamic data binding
- **Banking Standards**: Pre-configured templates for common financial documents
- **Security**: JWT authentication and static token verification
- **Scalable**: Docker/Kubernetes ready with CI/CD pipelines

## Getting Started

### Prerequisites

- Python 3.11+
- Docker 20.10+
- LibreOffice and Chromium (for PDF conversion)

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PDF_CONVERTER_URL` | PDF conversion service URL | `http://localhost:3100` |
| `STATIC_TOKEN` | Authentication token | `default_token` |
| `TEMPLATE_DIR` | Templates directory | `templates` |
| `SERVICE_CONTEXT_URL` | Base API path | `/document-generator` |

### API Reference

#### Major Endpoints

- `POST /api/v1/generate-docx` - Generate Word documents
- `POST /api/v1/generate-xlsx` - Generate Excel spreadsheets
- `POST /api/v1/generate-html` - Generate HTML documents
- `GET /api/v1/templates` - List available templates

### Example Request

```bash
curl -X POST "http://localhost:8000/api/v1/generate-docx" \
  -H "Authorization: Bearer your_token" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "ACCOUNT_STATEMENT",
    "format": "pdf",
    "data": {
      "clientName": "John Doe",
      "accountNumber": "KZ123456789",
      "transactions": [...]
    }
  }'
```

## Template Development

### HTML Templates
```docx
{{ clientName }}        <- Simple placeholder
{{#showLogo}}           <- Conditional block
  <img src="logo.jpg">
{{/showLogo}}

{{#table}}              <- Table row iteration
  <tr>
    <td>{{ date }}</td>
    <td>{{ amount }}</td>
  </tr>
{{/table}}
```

### XLSX Templates
```
{{ single_value }}      <- Cell value
{{ table.row_value }}   <- Table expansion
```

## Deployment

### Local Development
```bash
docker-compose up --build
```

## License

Proprietary software © Bank RBK. For internal use only.

## Шаблоны

## PDP - Заявление о досрочном погашении

- `code` — PDP
- `format` — формат генерируемого файла (**docx**, **pdf**).
- `data` — структура JSON метаданных должны быть следующей:
```json
{
  "clientName": "Имя или наименование клиента",
  "clientIin": "БИН/ИИН клиента",
  "sum": "Сумма для досрочного погашения",
  "accNumber": "Номер текущего счета клиента",
  "creditAgreementNumber": "Номер заявления на получение транша",
  "creditAgreementId": "ID кредитного договора",
  "creditAgreementDate": "Дата кредитного договора",
  "operDay": "Операционный день (дата подачи заявления)",
  "personName": "ФИО подписанта",
  "personPosition": "Должность подписанта"
}
```
## CURRENCY_CONTRACT - Заявление о принятии валютного договора по экспорту или импорту на валютный контроль

- `code` — CURRENCY_CONTRACT
- `format` — формат генерируемого файла (**docx**, **pdf**).
- `data` — структура JSON метаданных должны быть следующей:
```json
{
  "contractNumber": "Номер валютного договора",
  "contractDate": "Дата валютного договора",
  "contractCategory": "Категория валютного договора (экспорт/импорт)",
  "unkRegNumber": "Учетный номер валютного договора",
  "unkRegDate": "Дата присвоения учетного номера",
  "customerName": "Наименование экспортера или импортера",
  "bin": "БИН/ИИН экспортера или импортера",
  "address": "Адрес экспортера или импортера",
  "phone": "Телефон экспортера или импортера",
  "email": "Электронная почта экспортера или импортера",
  "beneficiaryName": "Наименование или ФИО иностранного партнера",
  "beneficiaryCountry": "Страна иностранного партнера",
  "contractAmount": "Ориентировочная сумма договора",
  "contractCurrency": "Валюта договора",
  "repatriationPeriod": "Срок репатриации",
  "contractTypesId": "Код вида валютного договора",
  "signature": "Подпись руководителя",
  "firstName": "Имя руководителя",
  "lastName": "Фамилия руководителя",
  "middleName": "Отчество руководителя",
  "createdDateTime": "Дата подачи заявления"
}
```

## REGISTRATION_APPLICATION_LLP - ЗАЯВЛЕНИЕ-АНКЕТА НА БАНКОВСКОЕ ОБСЛУЖИВАНИЕ В АО «BANK RBK»

- `code` — REGISTRATION_APPLICATION_LLP
- `format` — формат генерируемого файла (**docx**, **pdf**).
- `data` — структура JSON метаданных должны быть следующей:
```json
{
  "companyName": "Наименование юридического лица",
  "companyNameInLatin": "Наименование компании латиницей",
  "fullNameDirector": "ФИО уполномоченного лица",
  "docNo": "Номер удостоверяющего документа",
  "issuedBy": "Кем выдан документ",
  "issuedDate": "Дата выдачи документа",
  "bin": "БИН юридического лица",
  "okeds": "Код ОКЭД",
  "workReason": "Дата регистрации юр. лица",
  "phone": "Контактный телефон",
  "email": "Электронная почта",
  "addressFact": "Фактический адрес",
  "addressCorp": "Юридический адрес",
  "position": "Должность подписанта",
  "resident": true, // Признак резидентства
  "fatca": false, // по FATCA (Закон США «О налоговом контроле иностранных счетов» (Foreign Account Tax Compliance Act)
  "hasKztAccount": false, // Открытие счета в тенге
  "hasRubAccount": false, // Открытие счета в рублях
  "hasUsdAccount": false, // Открытие счета в долларах
  "hasEurAccount": false, // Открытие счета в евро
  "hasNotPredefinedCurrencyAccount": true, // Открытие счета в иной валюте
  "notPredefinedCurrencyName": "Кыргызский сом",
  "has10MillionWarrantyCompensation": true, // Гарантия до 10 млн тенге
  "has5MillionWarrantyCompensation": false,  // Гарантия до 5 млн тенге
  "day": "День подключения тарифа",
  "month": "Месяц подключения тарифа",
  "year": "Год подключения тарифа",
  "tariff": "Название тарифного пакета",
  "period": "Периодичность списания комиссии",
  "idn": "ИИН пользователя (если есть)",
  "birthDate": "Дата рождения",
  "docType": "Тип документа",
  "endDate": "Срок действия документа",
  "surname": "Фамилия участника/директора",
  "name": "Имя участника/директора",
  "patronymic": "Отчество участника/директора",
  "iin": "ИИН участника/директора",
  "registrationDepartment": "Орган регистрации документа",
  "currentDate": "Дата подачи заявления",
  "haveLicense": true, // Наличие лицензии
  "accNominalHolder": true, // Счет для учета денег Клиентов номинального держателя
  "accCustodian": true, // Счет для учета денег Клиентов кастодиана
  "commissionCurrencyTransaction": false, // Взимание комиссии В валюте проведения операции
  "commissionOtherCurrency": true, // Взимание комиссии В другой валюте
  "typePaymentCard":"VisaSignatureBusiness", // Вид платежной карточки VisaBusiness, VisaSignatureBusiness, Other
  "linkedNewCurrentAcc":true, // С привязкой к новому текущему счету Клиента
  "openCardAcc": true, //  С открытием отдельного Карт-счета(текущего)
  "internetClientOtp": true, // Прошу подключить к Системе «Интернет-Клиент»  OTP
  "internetClientSms": true // Прошу подключить к Системе «Интернет-Клиент»  SMS
}
```

## REGISTRATION_APPLICATION_IE - ЗАЯВЛЕНИЕ-АНКЕТА НА БАНКОВСКОЕ ОБСЛУЖИВАНИЕ В АО «BANK RBK»

- `code` — REGISTRATION_APPLICATION_IE
- `format` — формат генерируемого файла (**docx**, **pdf**).
- `data` — структура JSON метаданных должны быть следующей:
```json
{
  "companyName": "Наименование клиента (ИП)",
  "companyNameInLatin": "Наименование клиента латиницей",
  "fullNameDirector": "ФИО уполномоченного лица",
  "docNo": "Номер документа, удостоверяющего личность",
  "issuedBy": "Кем выдан документ",
  "issuedDate": "Дата выдачи документа",
  "bin": "БИН/ИИН клиента",
  "okeds": "Код ОКЭД",
  "workReason": "Дата регистрации ИП",
  "phone": "Контактный телефон",
  "email": "Электронная почта",
  "addressFact": "Фактический адрес",
  "addressCorp": "Юридический адрес",
  "position": "Должность уполномоченного лица",
  "resident": true, // Признак резидентства
  "fatca": false, // по FATCA (Закон США «О налоговом контроле иностранных счетов» (Foreign Account Tax Compliance Act)
  "hasKztAccount": false, // Открытие счета в тенге
  "hasRubAccount": false, // Открытие счета в рублях
  "hasUsdAccount": false, // Открытие счета в долларах
  "hasEurAccount": false, // Открытие счета в евро
  "hasNotPredefinedCurrencyAccount": true, // Открытие счета в иной валюте
  "notPredefinedCurrencyName": "Кыргызский сом",
  "has10MillionWarrantyCompensation": true, // Гарантия до 10 млн тенге
  "has5MillionWarrantyCompensation": false,  // Гарантия до 5 млн тенге
  "day": "День подключения тарифа",
  "month": "Месяц подключения тарифа",
  "year": "Год подключения тарифа",
  "tariff": "Название тарифного пакета",
  "period": "Периодичность списания комиссии",
  "idn": "ИИН пользователя",
  "birthDate": "Дата рождения",
  "docType": "Тип документа",
  "endDate": "Срок действия документа",
  "surname": "Фамилия ИП/представителя",
  "name": "Имя ИП/представителя",
  "patronymic": "Отчество ИП/представителя",
  "iin": "ИИН ИП/представителя",
  "registrationDepartment": "Орган регистрации",
  "currentDate": "Дата подачи заявления",
  "haveLicense": true, // Наличие лицензии
  "accNominalHolder": true, // Счет для учета денег Клиентов номинального держателя
  "accCustodian": true, // Счет для учета денег Клиентов кастодиана
  "commissionCurrencyTransaction": false, // Взимание комиссии В валюте проведения операции
  "commissionOtherCurrency": true, // Взимание комиссии В другой валюте
  "typePaymentCard":"VisaSignatureBusiness", // Вид платежной карточки VisaBusiness, VisaSignatureBusiness, Other
  "linkedNewCurrentAcc":true, // С привязкой к новому текущему счету Клиента
  "openCardAcc": true, //  С открытием отдельного Карт-счета(текущего)
  "internetClientOtp": true, // Прошу подключить к Системе «Интернет-Клиент»  OTP
  "internetClientSms": true // Прошу подключить к Системе «Интернет-Клиент»  SMS
}
```

## REGISTRATION_APPLICATION_FOR_CURRENT_USER - Заявление на подключение к Системе «Интернет-Клиент» (внесение изменений)

- `code` — REGISTRATION_APPLICATION_FOR_CURRENT_USER
- `format` — формат генерируемого файла (**docx**, **pdf**).
- `data` — структура JSON метаданных должны быть следующей:

```json
{
  "currentDate": "Дата подачи заявления",
  "companyName": "Юридическое наименование или ФИО клиента",
  "bin": "БИН/ИИН клиента",
  "resident": "Признак резидентства — резидент",
  "noResident": "Признак резидентства — нерезидент",
  "name": "Имя уполномоченного лица",
  "surname": "Фамилия уполномоченного лица",
  "patronymic": "Отчество уполномоченного лица",
  "docType": "Тип документа пользователя",
  "docNo": "Номер документа",
  "issuedDate": "Дата выдачи документа",
  "endDate": "Срок действия документа",
  "issuedBy": "Орган, выдавший документ",
  "phone": "Мобильный телефон пользователя",
  "email": "Электронная почта пользователя",
  "iin": "ИИН пользователя",
  "birthDate": "Дата рождения пользователя",
  "authType": "OTP" // OTP or SMS
}
```
## ACCOUNT_STATEMENT - Выписка по лицевому счету

- `code` — ACCOUNT_STATEMENT
- `format` — формат генерируемого файла (**docx**, **pdf**).
- `data` — структура JSON метаданных должны быть следующей:
```json
{
    "statementDate": "2025-04-11",
    "periodFrom": "2025-03-01",
    "periodTo": "2025-03-31",
    "ownerName": "Иванов Иван Иванович",
    "ownerBankName": "АО Банк Развития",
    "ownerBankBIK": "ABCDKZKX",
    "ownerTaxCode": "123456789012",
    "accNumber": "KZ123456789012345678",
    "address": "г. Алматы, ул. Абая, д. 10",
    "openingDate": "2022-01-15",
    "docType": "Текущий счет",
    "accCurrency": "KZT",
    "inAmount": "150000.00",
    "table": [
        {
            "reportCount": 1,
            "date": "2025-03-05",
            "docNumber": "INV-1001",
            "payer": "ТОО «Компания А»",
            "receiver": "Иванов Иван Иванович",
            "debet": "0.00",
            "credit": "50000.00",
            "paymentPurpose": "Оплата по договору №123 от 01.03.2025"
        },
        {
            "reportCount": 2,
            "date": "2025-03-20",
            "docNumber": "OUT-2001",
            "payer": "Иванов Иван Иванович",
            "receiver": "ТОО «Компания Б»",
            "debet": "30000.00",
            "credit": "0.00",
            "paymentPurpose": "Оплата услуг по счету №567 от 18.03.2025"
        }
    ],
    "debetTotal": "30000.00",
    "creditTotal": "50000.00",
    "documentCount": 2
}
```

## ACCOUNT_STATEMENT_PORTRAIT - Выписка по лицевому счету (портретная ориентация)

- `code` — ACCOUNT_STATEMENT_PORTRAIT
- `format` — формат генерируемого файла (**docx**, **pdf**).
- `data` — структура JSON метаданных должны быть следующей:
```json
{
    "statementDate": "2025-04-11",
    "periodFrom": "2025-03-01",
    "periodTo": "2025-03-31",
    "ownerName": "Иванов Иван Иванович",
    "ownerBankName": "АО Банк Развития",
    "ownerBankBIK": "ABCDKZKX",
    "ownerTaxCode": "123456789012",
    "accNumber": "KZ123456789012345678",
    "address": "г. Алматы, ул. Абая, д. 10",
    "openingDate": "2022-01-15",
    "docType": "Текущий счет",
    "accCurrency": "KZT",
    "inAmount": "150000.00",
    "table": [
        {
            "reportCount": 1,
            "date": "2025-03-05",
            "docNumber": "INV-1001",
            "payer": "ТОО «Компания А»",
            "receiver": "Иванов Иван Иванович",
            "debet": "0.00",
            "credit": "50000.00",
            "paymentPurpose": "Оплата по договору №123 от 01.03.2025"
        },
        {
            "reportCount": 2,
            "date": "2025-03-20",
            "docNumber": "OUT-2001",
            "payer": "Иванов Иван Иванович",
            "receiver": "ТОО «Компания Б»",
            "debet": "30000.00",
            "credit": "0.00",
            "paymentPurpose": "Оплата услуг по счету №567 от 18.03.2025"
        }
    ],
    "debetTotal": "30000.00",
    "creditTotal": "50000.00",
    "documentCount": 2
}
```

## CREDIT_SCHEDULE - График платежей
- `code` — CREDIT_SCHEDULE
- `format` — формат генерируемого файла (**docx**, **pdf**).
- `data` — структура JSON метаданных должны быть следующей:
```json
{
    "documentNumber": "12345-КД-2025",
    "reportCreationDate": "2025-04-11",
    "clientFullName": "ТОО «Пример Компания»",
    "taxCode": "123456789012",
    "table": [
        {
            "scheduleDate": "2025-05-01",
            "amountOfLoanToMaturity": "500000.00",
            "amountOfRemunerToMaturity": "25000.00",
            "principalBalance": "4500000.00",
            "amountOfForthcomingPayment": "525000.00"
        },
        {
            "scheduleDate": "2025-06-01",
            "amountOfLoanToMaturity": "500000.00",
            "amountOfRemunerToMaturity": "22000.00",
            "principalBalance": "4000000.00",
            "amountOfForthcomingPayment": "522000.00"
        }
    ]
}
```

## CREDIT_SCHEDULE_WITH_DAMU - График платежей с процентами ДАМУ

- `code` — CREDIT_SCHEDULE_WITH_DAMU
- `format` — формат генерируемого файла (**docx**, **pdf**).
- `data` — структура JSON метаданных должны быть следующей:
```json
{
    "documentNumber": "12345-КД-2025",
    "reportCreationDate": "2025-04-11",
    "clientFullName": "ТОО «Пример Компания»",
    "taxCode": "123456789012",
    "table": [
        {
            "scheduleDate": "2025-05-01",
            "amountOfLoanToMaturity": "500000.00",
            "amountOfRemunerToMaturity": "25000.00",
            "principalBalance": "4500000.00",
            "amountOfForthcomingPayment": "525000.00",
            "amountTaxDamu": "10000.00"
        },
        {
            "scheduleDate": "2025-06-01",
            "amountOfLoanToMaturity": "500000.00",
            "amountOfRemunerToMaturity": "22000.00",
            "principalBalance": "4000000.00",
            "amountOfForthcomingPayment": "522000.00",
            "amountTaxDamu": "10000.00"
        }
    ]
}
```

## PAYMENT_ORDER - Платежное поручение

- `code` — PAYMENT_ORDER
- `format` — формат генерируемого файла (**docx**, **pdf**).
- `data` — структура JSON метаданных должны быть следующей:
```json
{
    "docNumber": "001234",
    "created": "2025-04-11",
    "ultimateDebtorName": "ТОО «Компания А»",
    "payerName": "ТОО «Компания А»",
    "payerTax": "123456789012",
    "PayerResidencyAndEconomicCode": "RZ001",
    "amount": "1000000.00",
    "payerCountry": "Казахстан",
    "payerBank": "АО «Банк Отправителя»",
    "payerBankCode": "BANKKZKX",
    "ultimateCreditorName": "ТОО «Компания Б»",
    "benefName": "ТОО «Компания Б»",
    "benefTax": "987654321098",
    "benefAcc": "KZ75999KZT1234567890",
    "benefRedideAndEconomicCode": "NR001",
    "benefCountry": "Казахстан",
    "benefBankName": "АО «Банк Получателя»",
    "benefBankCode": "BENFKZKX",
    "amountText": "Один миллион тенге",
    "paymentPurpose": "Оплата за поставку оборудования по договору №45 от 01.04.2025",
    "paymentPurposeCode": "771",
    "budgetClassificationCode": "302101",
    "valueDate": "2025-04-12",
    "director": "Иванов Иван Иванович",
    "accountant": "Петров Петр Петрович"
}
```

## CARD_STATEMENT_EXCEL - Выписка по карте(xlsx)

- `code` — CARD_STATEMENT
- `format` — формат генерируемого файла (**xlsx**).
- `data` — структура JSON метаданных должны быть следующей:
```json
{
  "date": "Дата формирования документа",
  "accountNumber": "Номер счета",
  "accountCurrency": "Валюта счета",
  "cardName": "Название карты",
  "clientName": "Имя клиента",
  "clientTaxCode": "ИИН/БИН клиента",
  "statementDateFrom": "Начало периода выписки",
  "statementDateTo": "Конец периода выписки",
  "currentTime": "Текущая дата/время",
  "blockedSum": "Заблокированная сумма",
  "initialBalance": "Начальный остаток",
  "income": "Сумма поступлений",
  "expenses": "Сумма списаний",
  "finalBalance": "Конечный остаток",
  "cardHoldLabel": "Пометка держателя карты",
  "statement": [
    {
      "transactionDate": "Дата операции",
      "processDate": "Дата обработки",
      "description": "Описание операции",
      "amount": "Сумма операции в валюте операции",
      "cardAmount": "Сумма операции в валюте карты",
      "commission": "Комиссия по операции"
    }
  ]
}
```

## ACCOUNT_STATEMENT_EXCEL - Выписка по лицевому счету(xlsx)

- `code` — ACCOUNT_STATEMENT_EXCEL
- `format` — формат генерируемого файла (**xlsx**, **pdf**).
- `data` — структура JSON метаданных должны быть следующей:
```json
{
  "periodFrom": "Дата начала периода выписки",
  "periodTo": "Дата окончания периода выписки",
  "ownerName": "Наименование клиента",
  "ownerBankName": "Наименование банка клиента",
  "ownerBankBik": "БИК банка клиента",
  "ownerTaxCode": "БИН/ИИН клиента",
  "accNumber": "Номер счета",
  "address": "Адрес клиента",
  "openingDate": "Дата открытия счета",
  "docType": "Тип документа",
  "accCurrency": "Валюта счета",
  "inAmount": "Начальный остаток",
  "outAmount": "Конечный остаток",
  "creditTotal": "Сумма по кредиту",
  "debetTotal": "Сумма по дебету",
  "documentCount": "Количество документов",
  "statement": [
    {
      "no": "Номер строки",
      "date": "Дата документа",
      "sender": "Отправитель (ФИО/наименование или реквизиты)",
      "receiver": "Получатель (ФИО/наименование или реквизиты)",
      "debet": "Сумма по дебету",
      "credit": "Сумма по кредиту",
      "paymentPurpose": "Назначение платежа"
    }
  ]
}
```

## EXTENDED_STATEMENT_EXCEL - Выписка по движениям средств(xlsx)

- `code` — EXTENDED_STATEMENT_EXCEL
- `format` — формат генерируемого файла (**xlsx**).
- `data` — структура JSON метаданных должны быть следующей:
```json
{
  "periodFrom": "Дата начала периода выписки",
  "periodTo": "Дата окончания периода выписки",
  "statementDate": "Дата формирования документа",
  "ownerName": "Наименование клиента",
  "ownerTaxCode": "БИН/ИИН клиента",
  "address": "Адрес клиента",
  "accNumber": "Номер счета",
  "openingDate": "Дата открытия счета",
  "accCurrency": "Валюта счета",
  "ownerBankName": "Наименование банка",
  "ownerBankBIK": "БИК банка",
  "inAmount": "Начальный остаток",
  "debetTotal": "Итого по дебету",
  "creditTotal": "Итого по кредиту",
  "statement": [
    {
      "no": "Номер строки",
      "date": "Дата документа",
      "docNumber": "Номер документа",
      "docType": "Тип документа",
      "purposeCode": "Код назначения платежа (КНП)",
      "beneficiaryKbe": "КБе получателя",
      "beneficiaryName": "Наименование получателя",
      "beneficiaryIban": "Счет получателя (IBAN)",
      "beneficiaryBin": "БИН получателя",
      "beneficiaryBankName": "Наименование банка получателя",
      "beneficiaryBankBic": "БИК банка получателя",
      "valueDate": "Дата валютирования",
      "debet": "Сумма по дебету",
      "credit": "Сумма по кредиту",
      "paymentPurpose": "Назначение платежа",
      "exchange": "Курс обмена",
      "amountNational": "Сумма в нац. валюте"
    }
  ]
}
```

## BUSINESS_ACCOUNT_STATEMENT_EXCEL - Выписка по лицевому счету(xlsx)

- `code` — BUSINESS_ACCOUNT_STATEMENT_EXCEL
- `format` — формат генерируемого файла (**xlsx**, **pdf**).
- `data` — структура JSON метаданных должны быть следующей:
```json
{
  "periodFrom": "Дата начала периода выписки",
  "periodTo": "Дата окончания периода выписки",
  "ownerName": "Наименование клиента",
  "ownerBankName": "Наименование банка клиента",
  "ownerBankBIK": "БИК банка клиента",
  "ownerTaxCode": "БИН/ИИН клиента",
  "accNumber": "Номер счета (ИИК)",
  "address": "Адрес клиента",
  "openingDate": "Дата открытия счета",
  "accCurrency": "Валюта счета",
  "inAmount": "Начальный остаток",
  "outAmount": "Конечный остаток",
  "creditTotal": "Общий кредитовый оборот",
  "debetTotal": "Общий дебетовый оборот",
  "documentCount": "Количество документов",
  "statement": [
    {
      "no": "Номер строки",
      "transactionDate": "Дата операции",
      "processDate": "Дата обработки",
      "docNumber": "Номер документа",
      "requisites": "Реквизиты корреспондента",
      "debet": "Сумма по дебету",
      "credit": "Сумма по кредиту",
      "commission": "Комиссия",
      "description": "Назначение платежа",
      "knp": "Код назначения платежа (КНП)"
    }
  ]
}
```

## SWIFT_INFO - информация по SWIFT-платежам

- `code` — SWIFT_INFO
- `format` — формат генерируемого файла (**docx**, **pdf**).
- `data` — структура JSON метаданных должны быть следующей:
```json
{
  "swift": "SWIFT"
}
```

## CONVERTATION - Конвертация документов
- `code` — CONVERTATION
- `format` — формат генерируемого файла (**docx**, , **pdf**).
- `data` — структура JSON метаданных должны быть следующей:
```json
{
  "docNumber": "001234",
  "created": "2025-04-11",
  "payerName": "ТОО «Компания А»",
  "payerTax": "123456789012",
  "payerResidencyAndEconomicCode1": "RZ001",
  "payerResidencyAndEconomicCode2": "RZ001",
  "valueDate": "2025-04-12",
  "recieverResidencyAndEconomicCode1": "NR001",
  "recieverResidencyAndEconomicCode2": "NR001",
  "firstCurrency": "USD",
  "accountFirst": "KZ123456789012345678",
  "secondCurrency": "EUR", 
  "accountSecond": "KZ987654321098765432",
  "amountFirst": "1000.00",
  "amountTextFirst": "Одна тысяча долларов США",
  "amountSecond": "900.00",
  "amountTextSecond": "Девятьсот евро",
  "codeKnp": "771",
  "rate": "1.11",
  "operationTargetCodeAndName": "Конвертация валюты"
}
```

## EXTENDED_STATEMENT - Выписка по движениям средств

---
### ⚠️ Дисклеймер

Поля `exchange` и `amountNational` присутствуют в структуре `table`, а также в шаблоне документа (отображаются в правой части таблицы).  
Однако, **в финальную выписку они не попадают**. Если потребуется их отобразить, необходимо обновить шаблон.
---

- `code` — EXTENDED_STATEMENT
- `format` — формат генерируемого файла (**docx**, **pdf**).
- `data` — структура JSON метаданных должны быть следующей:'

```json
{
  "periodFrom": "Дата начала периода выписки",
  "periodTo": "Дата окончания периода выписки",
  "statementDate": "Дата формирования документа",
  "ownerName": "Наименование клиента",
  "ownerTaxCode": "БИН/ИИН клиента",
  "address": "Адрес клиента",
  "accNumber": "Номер счета",
  "openingDate": "Дата открытия счета",
  "accCurrency": "Валюта счета",
  "ownerBankName": "Наименование банка",
  "ownerBankBik": "БИК банка",
  "inAmount": "Начальный остаток",
  "outAmount": "Конечный остаток",
  "documentCount": "Количество документов",
  "debetTotal": "Итого по дебету",
  "creditTotal": "Итого по кредиту",
  "table": [
    {
      "reportCount": "Номер строки",
      "date": "Дата документа",
      "docNumber": "Номер документа",
      "docType": "Тип документа",
      "purposeCode": "Код назначения платежа (КНП)",
      "beneficiaryKbe": "КБе получателя",
      "beneficiaryName": "Наименование получателя",
      "beneficiaryIban": "Счет получателя (IBAN)",
      "beneficiaryBin": "БИН получателя",
      "beneficiaryBankName": "Наименование банка получателя",
      "beneficiaryBankBic": "БИК банка получателя",
      "valueDate": "Дата валютирования",
      "debet": "Сумма по дебету",
      "credit": "Сумма по кредиту",
      "paymentPurpose": "Назначение платежа",
      
      "exchange": "Курс обмена",
      "amountNational": "Сумма в нац. валюте"
    }
  ]
}
```

## BUSINESS_ACCOUNT_STATEMENT - Выписка по лицевому счету
- `code` — BUSINESS_ACCOUNT_STATEMENT
- `format` — формат генерируемого файла (**docx**, **pdf**).
- `data` — структура JSON метаданных должны быть следующей:
```json
{
  "statementDate": "Дата формирования документа",
  "periodFrom": "Дата начала периода выписки",
  "periodTo": "Дата окончания периода выписки",
  "ownerName": "Наименование клиента",
  "ownerBankName": "Наименование банка клиента",
  "ownerBankBIK": "БИК банка клиента",
  "ownerTaxCode": "БИН/ИИН клиента",
  "accNumber": "Номер счета (ИИК)",
  "address": "Адрес клиента",
  "openingDate": "Дата открытия счета",
  "accCurrency": "Валюта счета",
  "inAmount": "Начальный остаток",
  "outAmount": "Конечный остаток",
  "availableAmount": "Доступный остаток",
  "documentCount": "Количество документов",
  "creditTotal": "Общий кредитовый оборот",
  "debetTotal": "Общий дебетовый оборот",
  "table": [
    {
      "reportCount": "Номер строки",
      "opDate": "Дата операции",
      "docNumber": "Номер документа",
      "requisites": "Реквизиты корреспондента",
      "debet": "Сумма по дебету",
      "credit": "Сумма по кредиту",
      "commission": "Комиссия",
      "description": "Назначение платежа",
      "knp": "Код назначения платежа (КНП)",
      "codeName": "Наименование кода"
    }
  ]
}
```

## CARD_STATEMENT - Выписка по карте
- `code` — CARD_STATEMENT
- `format` — формат генерируемого файла (**docx**, **pdf**).
- `data` — структура JSON метаданных должны быть следующей:
```json
{
  "date": "Дата формирования документа",
  "accountNumber": "Номер счета",
  "accountCurrency": "Валюта счета",
  "clientName": "Имя клиента",
  "clientTaxCode": "ИИН/БИН клиента",
  "statementDateFrom": "Начало периода выписки",
  "statementDateTo": "Конец периода выписки",
  "currentTime": "Текущая дата/время",
  "blockedSum": "Заблокированная сумма",
  "initialBalance": "Начальный остаток",
  "income": "Сумма поступлений",
  "expenses": "Сумма списаний",
  "finalBalance": "Конечный остаток",
  "table1": [
    {
      "dCreationTime": "Дата операции",
      "dProcessingDate": "Дата обработки",
      "dDescription": "Описание операции",
      "dOperationAmount": "Сумма операции в валюте операции",
      "dAccountAmount": "Сумма операции в валюте карты",
      "dCommission": "Комиссия по операции"
    }
  ],
  table2: [
    {
      "creationTime": "Дата операции",
      "processingDate": "Дата обработки",
      "description": "Описание операции",
      "operationAmount": "Сумма операции в валюте операции",
      "accountAmount": "Сумма операции в валюте карты"
    }
  ]
}
```
{
  "code": "CARD_STATEMENT_EXCEL",
  "format": "xlsx",
  "data": {
    "date": "Дата формирования документа",
    "accountNumber": "Номер счета",
    "accountCurrency": "Валюта счета",
    "cardName": "Kaspi gold",
    "clientName": "Имя клиента",
    "clientTaxCode": "ИИН/БИН клиента",
    "statementDateFrom": "Начало периода выписки",
    "statementDateTo": "Конец периода выписки",
    "currentTime": "Текущая дата/время",
    "blockedSum": "Заблокированная сумма",
    "initialBalance": "Начальный остаток",
    "income": "Сумма поступлений",
    "expenses": "Сумма списаний",
    "finalBalance": "Конечный остаток",
    "transactions": [
      {
        "creationTime": "Дата операции",
        "processingTime": "Дата обработки",
        "description": "Описание операции",
        "operationAmount": "Сумма операции в валюте операции",
        "accountAmount": "Сумма операции в валюте карты",
        "commission": "Комиссия по операции"
      }
    ],
    "waitTransactions": [
      {
        "creationTime": "Дата операции",
        "processingTime": "Дата обработки",
        "description": "Описание операции",
        "operationAmount": "Сумма операции в валюте операции",
        "accountAmount": "Сумма операции в валюте карты"
      }
    ]
  }
}
